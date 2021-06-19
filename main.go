package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"sync"
	"time"

	"github.com/gocql/gocql"
	"github.com/webdevgopi/chatApp-gRPC/db"
	"github.com/webdevgopi/chatApp-gRPC/interceptors"
	"github.com/webdevgopi/chatApp-gRPC/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
)

type StreamConn struct {
	error  chan error
	stream proto.MessagingService_ConnectUserServer
	status proto.User_UserStatusCode
	device string
}

type Connection struct {
	streamConnections []*StreamConn
	userObj           *proto.User
}

type Server struct {
	OnlineUsers map[string]*Connection
	proto.UnimplementedMessagingServiceServer
}

func getUserObjFromDB(userId string) (*proto.User, error) {
	var name string
	var uType string
	if err := session.Query(`SELECT name, user_type FROM chat_app.users WHERE user_id = ? `,
		userId).Scan(&name, &uType); err != nil {
		if err == gocql.ErrNotFound {
			return nil, status.Error(codes.InvalidArgument, "there's no user with that UUID")
		}
		return nil, status.Error(codes.Internal, "system is broken")
	}
	var userType proto.User_Type
	switch uType {
	case "expert":
		userType = proto.User_EXPERT
	case "patient":
		userType = proto.User_PATIENT
	default:
		userType = proto.User_UNKNOWN
	}
	return &proto.User{
		UserId: userId,
		Name:   name,
		Status: proto.User_LOGGED_OUT,
		Type:   userType,
	}, nil
}

func addToUserConnObj(s Server, server proto.MessagingService_ConnectUserServer, user *proto.ConnectUserRequest) {
	conn, ok := s.OnlineUsers[user.GetUserId()]
	if !ok {
		streamConns := make([]*StreamConn, 0, 1)
		streamConns = append(streamConns, &StreamConn{
			stream: server,
			error:  make(chan error),
			status: proto.User_ONLINE,
			device: user.GetDevice(),
		})
		userType := proto.User_UNKNOWN
		if user.GetType() == proto.ConnectUserRequest_EXPERT {
			userType = proto.User_EXPERT
		}
		if user.GetType() == proto.ConnectUserRequest_PATIENT {
			userType = proto.User_PATIENT
		}
		s.OnlineUsers[user.GetUserId()] = &Connection{
			streamConnections: streamConns,
			userObj: &proto.User{
				UserId: user.GetUserId(),
				Name:   user.GetName(),
				Status: proto.User_ONLINE,
				Type:   userType,
			},
		}
	} else {
		conn.streamConnections = append(conn.streamConnections, &StreamConn{
			stream: server,
			error:  make(chan error),
			status: proto.User_ONLINE,
			device: user.GetDevice(),
		})
	}
}

func getUserObj(s Server, userId string) (*proto.User, error) {
	elem, ok := getUserConnObj(s, &proto.User{
		UserId: userId,
	})
	if !ok {
		user, err := getUserObjFromDB(userId)
		if err != nil {
			return nil, err
		}
		return user, nil
	}
	return elem.userObj, nil

}

func getUserConnObj(s Server, user *proto.User) (*Connection, bool) {
	elem, ok := s.OnlineUsers[user.GetUserId()]
	if ok {
		if len(elem.streamConnections) == 0 {
			delete(s.OnlineUsers, user.GetUserId())
			return nil, false
		}
		userStatus := proto.User_OFFLINE
		for _, conn := range elem.streamConnections {
			if conn.status == proto.User_ONLINE {
				userStatus = conn.status
			}
		}
		elem.userObj.Status = userStatus
		return elem, true
	}
	return nil, false
}

func getUserUnreadMsgs(user *proto.User, s Server) []*proto.ChatMessage {

	queryResults := make([]*proto.ChatMessage, 0, 6)
	msgStatus1 := "delivered"
	msgStatus2 := "notDelivered"
	msgStatus3 := "failed"
	scan := session.Query(`SELECT chat_id, body, time, from_user, reply_for_chat_id FROM chat_app.chat_messages WHERE to_user = ? AND status = ? ALLOW FILTERING `, user.GetUserId(), msgStatus1).Iter().Scanner()
	for scan.Next() {
		var (
			body           string
			chatId         string
			msgTime        time.Time
			fromUser       string
			replyForChatId string
		)
		err2 := scan.Scan(&chatId, &body, &msgTime, &fromUser, &replyForChatId)
		if err2 != nil {
			if err2 == gocql.ErrNotFound {
				break
			}
			log.Printf("error while querying; %v", err2.Error())
		}

		b, err2 := msgTime.MarshalText()
		if err2 != nil {
			log.Printf("error while marshaling time; %v\n", err2.Error())
		}
		u, _ := getUserObj(s, fromUser)

		queryResults = append(queryResults, &proto.ChatMessage{
			ChatId:         chatId,
			Body:           body,
			TimeStamp:      string(b),
			FromUser:       u,
			ToUser:         user,
			ReplyForChatId: replyForChatId,
		})
	}
	scan2 := session.Query(`SELECT chat_id, body, time, from_user, reply_for_chat_id FROM chat_app.chat_messages WHERE to_user = ? AND status = ? ALLOW FILTERING `, user.GetUserId(), msgStatus2).Iter().Scanner()
	for scan2.Next() {
		var (
			body           string
			chatId         string
			msgTime        time.Time
			fromUser       string
			replyForChatId string
		)
		err2 := scan2.Scan(&chatId, &body, &msgTime, &fromUser)
		if err2 != nil {
			if err2 == gocql.ErrNotFound {
				break
			}
			log.Printf("error while querying; %v", err2.Error())
		}

		b, err2 := msgTime.MarshalText()
		if err2 != nil {
			log.Printf("error while marshaling time; %v\n", err2.Error())
		}
		u, _ := getUserObj(s, fromUser)

		queryResults = append(queryResults, &proto.ChatMessage{
			ChatId:         chatId,
			Body:           body,
			TimeStamp:      string(b),
			FromUser:       u,
			ToUser:         user,
			ReplyForChatId: replyForChatId,
		})
	}
	scan3 := session.Query(`SELECT chat_id, body, time, from_user, reply_for_chat_id FROM chat_app.chat_messages WHERE to_user = ? AND status = ? ALLOW FILTERING `, user.GetUserId(), msgStatus3).Iter().Scanner()
	for scan3.Next() {
		var (
			body           string
			chatId         string
			msgTime        time.Time
			fromUser       string
			replyForChatId string
		)
		err2 := scan3.Scan(&chatId, &body, &msgTime, &fromUser, &replyForChatId)
		if err2 != nil {
			if err2 == gocql.ErrNotFound {
				break
			}
			log.Printf("error while querying; %v", err2.Error())
		}

		b, err2 := msgTime.MarshalText()
		if err2 != nil {
			log.Printf("error while marshaling time; %v\n", err2.Error())
		}
		u, _ := getUserObj(s, fromUser)

		queryResults = append(queryResults, &proto.ChatMessage{
			ChatId:         chatId,
			Body:           body,
			TimeStamp:      string(b),
			FromUser:       u,
			ToUser:         user,
			ReplyForChatId: replyForChatId,
		})
	}

	return queryResults
}

func getUserActiveConns(user *proto.User, s Server) ([]*StreamConn, error) {
	elem, ok := s.OnlineUsers[user.GetUserId()]
	if !ok {
		return nil, fmt.Errorf("no active connections to the given user")
	}
	count := 0
	for i, conn := range elem.streamConnections {
		if conn.status == proto.User_LOGGED_OUT {
			count++
			elem.streamConnections[i] = elem.streamConnections[len(elem.streamConnections)-count]
		}
	}
	elem.streamConnections = elem.streamConnections[0 : len(elem.streamConnections)-count]
	if len(elem.streamConnections) == 0 {
		delete(s.OnlineUsers, user.GetUserId())
		return nil, nil
	}
	return elem.streamConnections, nil
}

func (s Server) GetUser(ctx context.Context, user *proto.QueryUser) (*proto.User, error) {
	//log.Printf("Is authenticated ? : %v\n", ctx.Value("isAuthenticated"))

	if err := session.Query(`INSERT INTO chat_app.users (user_id, name, user_type, last_update_timestamp) VALUES (?,?,?,?) IF NOT EXISTS;`, user.GetUserId(), fmt.Sprintf("user %v", ctx.Value("user_uuid")), "expert", time.Now()).Exec(); err != nil {
		log.Fatal(err)
	}

	return getUserObj(s, user.GetUserId())
}

func (s Server) ConnectUser(request *proto.ConnectUserRequest, server proto.MessagingService_ConnectUserServer) error {
	//log.Printf("Is authenticated ? : %v\n", server.Context().Value("isAuthenticated"))
	addToUserConnObj(s, server, request)
	user, _ := getUserObj(s, request.GetUserId())
	unReadMsgs := getUserUnreadMsgs(user, s)
	conns, _ := getUserActiveConns(user, s)
	streamObj := conns[len(conns)-1]
	if len(unReadMsgs) > 0 {

		go func(msgs []*proto.ChatMessage, streamConn *StreamConn) {

			for _, msg := range unReadMsgs {
				err := server.Send(msg)
				if err != nil {
					streamConn.status = proto.User_LOGGED_OUT
					streamConn.error <- err
					return
				}
			}
		}(unReadMsgs, streamObj)

	}
	return <-streamObj.error
}

func (s Server) DisconnectUser(ctx context.Context, req *proto.DisconnectUserRequest) (*proto.User, error) {
	//log.Printf("Is authenticated ? : %v\n", ctx.Value("isAuthenticated"))
	user, err := getUserObj(s, req.GetUser().GetUserId())
	if err != nil {
		return nil, err
	}
	conns, err := getUserActiveConns(user, s)
	if err != nil {
		return nil, err
	}
	for _, conn := range conns {
		if conn.device == req.GetDevice() {
			switch sts := req.GetUser().GetStatus(); {
			case sts == proto.User_OFFLINE:
				conn.status = sts
			case sts == proto.User_LOGGED_OUT:
				conn.status = sts
				conn.error <- errors.New("EOF")
			default:
				conn.status = proto.User_UNKNOWN_STATUS
			}
		}
	}
	return getUserObj(s, req.GetUser().GetUserId())
}

func (s Server) BroadcastMessage(ctx context.Context, message *proto.ChatMessage) (*proto.ChatMessageResponse, error) {
	//log.Printf("Is authenticated ? : %v\n", ctx.Value("isAuthenticated"))
	wait := sync.WaitGroup{}
	done := make(chan int)
	_, ok := getUserConnObj(s, message.GetToUser())
	if ok {
		activeConns, _ := getUserActiveConns(message.GetToUser(), s)
		//log.Printf("activeConns %v \n", len(activeConns))
		msgStatuses := make(chan proto.ChatMessageResponse_MessageStatus, len(activeConns))
		for _, conn := range activeConns {
			wait.Add(1)
			go func(msg *proto.ChatMessage, conn *StreamConn, msgStatuses chan proto.ChatMessageResponse_MessageStatus) {
				defer wait.Done()
				err := conn.stream.Send(msg)
				//grpcLog.Info("Sending message to: ", conn.stream)
				if err != nil {
					//grpcLog.Errorf("Error with Stream: %v - Error: %v", conn.stream, err)
					conn.status = proto.User_LOGGED_OUT
					msgStatuses <- proto.ChatMessageResponse_FAILED
					conn.error <- err
					return
				} else {
					msgStatuses <- proto.ChatMessageResponse_DELIVERED
				}
			}(message, conn, msgStatuses)
		}

		go func() {
			wait.Wait()
			close(done)
			close(msgStatuses)
		}()

		<-done
		msgStatus := proto.ChatMessageResponse_FAILED
		msgStatusDB := "failed"
		for m := range msgStatuses {
			if m == proto.ChatMessageResponse_DELIVERED {
				msgStatus = m
				msgStatusDB = "delivered"
			}
		}

		var msgTime time.Time
		if err := msgTime.UnmarshalText([]byte(message.GetTimeStamp())); err != nil {
			log.Printf("Trouble while parsing Time; %v\n", err.Error())
		}
		if err := session.Query(`INSERT INTO chat_app.chat_messages (chat_id, from_user, to_user, body, status, time, reply_for_chat_id) VALUES (?,?,?,?,?,?,?)`, message.GetChatId(), message.GetFromUser().GetUserId(), message.GetToUser().GetUserId(), message.GetBody(), msgStatusDB, msgTime, message.GetReplyForChatId()).Exec(); err != nil {
			log.Fatal(err)
		}
		return &proto.ChatMessageResponse{
			ChatMessage: &proto.ChatMessage{
				ChatId: message.GetChatId(),
			},
			MsgStatus: msgStatus,
		}, nil
	}
	var msgTime time.Time
	if err := msgTime.UnmarshalText([]byte(message.GetTimeStamp())); err != nil {
		log.Printf("Trouble while parsing Time; %v\n", err.Error())
	}
	if err := session.Query(`INSERT INTO chat_app.chat_messages (chat_id, from_user, to_user, body, status, time, reply_for_chat_id) VALUES (?,?,?,?,?,?,?)`, message.GetChatId(), message.GetFromUser().GetUserId(), message.GetToUser().GetUserId(), message.GetBody(), "notDelivered", msgTime, message.GetReplyForChatId()).Exec(); err != nil {
		log.Fatal(err)
	}
	return &proto.ChatMessageResponse{
		ChatMessage: &proto.ChatMessage{
			ChatId: message.GetChatId(),
		},
		MsgStatus: proto.ChatMessageResponse_NOT_DELIVERED,
	}, nil
}

func (s Server) UpdateMessageStatus(ctx context.Context, chatMessageStatus *proto.ChatMessageStatus) (*proto.ChatMessageStatus, error) {
	//log.Printf("Is authenticated ? : %v\n", ctx.Value("isAuthenticated"))
	msg := chatMessageStatus.GetChatMessage()
	var msgStatus string
	switch sts := chatMessageStatus.GetMsgStatus(); {
	case sts == proto.ChatMessageStatus_DELIVERED:
		msgStatus = "delivered"
	case sts == proto.ChatMessageStatus_FAILED:
		msgStatus = "failed"
	case sts == proto.ChatMessageStatus_NOT_DELIVERED:
		msgStatus = "notDelivered"
	case sts == proto.ChatMessageStatus_SEEN:
		msgStatus = "seen"
	default:
		msgStatus = "unknown"
	}

	var msgTime time.Time
	if err := msgTime.UnmarshalText([]byte(msg.GetTimeStamp())); err != nil {
		log.Printf("Trouble while parsing Time; %v\n", err.Error())
		return nil, status.Error(codes.InvalidArgument, "invalid value for msg time property")
	}

	if err := session.Query(`UPDATE chat_app.chat_messages SET status = ? WHERE from_user = ? AND to_user = ? AND time = ?`, msgStatus, msg.GetFromUser().GetUserId(), msg.GetToUser().GetUserId(), msgTime).Exec(); err != nil {
		log.Printf("Trouble while updating db; %v\n", err.Error())
		return nil, status.Error(codes.Internal, "trouble updating message status")
	}

	return &proto.ChatMessageStatus{
		MsgStatus:   chatMessageStatus.GetMsgStatus(),
		ChatMessage: msg,
	}, nil
}

func (s Server) GetMessageStatus(ctx context.Context, message *proto.ChatMessage) (*proto.ChatMessageStatus, error) {
	//log.Printf("Is authenticated ? : %v\n", ctx.Value("isAuthenticated"))
	var msgTime time.Time
	if err := msgTime.UnmarshalText([]byte(message.GetTimeStamp())); err != nil {
		log.Printf("Trouble while parsing Time; %v\n", err.Error())
		return nil, status.Error(codes.InvalidArgument, "invalid value for msg time property")
	}
	var statusDB string
	if err := session.Query(`SELECT status FROM chat_app.chat_messages WHERE to_user = ? AND from_user = ? AND time = ?`, message.GetToUser().GetUserId(), message.GetFromUser().GetUserId(), msgTime).
		Scan(&statusDB); err != nil {
		if err == gocql.ErrNotFound {
			return nil, status.Error(codes.InvalidArgument, "invalid query parameters")
		}
		log.Printf("Trouble querying msg status; %v\n", err.Error())
		return nil, status.Error(codes.Internal, "trouble while querying msg data from")
	}
	var msgStatus proto.ChatMessageStatus_MessageStatus
	switch sts := statusDB; {
	case sts == "delivered":
		msgStatus = proto.ChatMessageStatus_DELIVERED
	case sts == "notDelivered":
		msgStatus = proto.ChatMessageStatus_NOT_DELIVERED
	case sts == "failed":
		msgStatus = proto.ChatMessageStatus_FAILED
	case sts == "seen":
		msgStatus = proto.ChatMessageStatus_SEEN
	}
	return &proto.ChatMessageStatus{ChatMessage: message, MsgStatus: msgStatus}, nil

}

func (s Server) GetMessages(ctx context.Context, request *proto.QueryMessagesRequest) (*proto.QueryMessagesResponse, error) {
	//log.Printf("Is authenticated ? : %v\n", ctx.Value("isAuthenticated"))
	var timeConstraint time.Time
	if err := timeConstraint.UnmarshalText([]byte(request.GetTimeConstraint())); err != nil {
		log.Printf("Trouble while parsing Time; %v\n", err.Error())
		return nil, status.Error(codes.InvalidArgument, "invalid value for msg time property")
	}
	var chatMessages []*proto.ChatMessage
	//log.Printf("querying messages\n")
	scan := session.Query(`SELECT chat_id, body, time, from_user ,to_user FROM chat_app.chat_messages WHERE to_user IN ( ?, ?) AND from_user IN ( ?, ? ) AND time <= ? LIMIT ? ALLOW FILTERING `, request.GetUser1().GetUserId(), request.GetUser2().GetUserId(), request.GetUser2().GetUserId(), request.GetUser1().GetUserId(), timeConstraint, request.GetLimit()).Iter().Scanner()
	for scan.Next() {
		var (
			body     string
			chatId   string
			msgTime  time.Time
			fromUser string
			toUser   string
		)
		err2 := scan.Scan(&chatId, &body, &msgTime, &fromUser, &toUser)
		if err2 != nil {
			log.Printf("error while querying; %v", err2.Error())
			return nil, status.Error(codes.Internal, "trouble while querying messages")
		}
		b, err2 := msgTime.MarshalText()
		if err2 != nil {
			log.Printf("error while marshaling time; %v\n", err2.Error())
			return nil, status.Error(codes.Internal, "trouble while querying messages")
		}
		var fromUserObj *proto.User
		var toUserObj *proto.User
		if fromUser == request.GetUser1().GetUserId() {
			fromUserObj = request.GetUser1()
			toUserObj = request.GetUser2()
		} else {
			fromUserObj = request.GetUser2()
			toUserObj = request.GetUser1()
		}
		chatMessages = append(chatMessages, &proto.ChatMessage{
			ChatId:    chatId,
			FromUser:  fromUserObj,
			ToUser:    toUserObj,
			Body:      body,
			TimeStamp: string(b),
		})
	}
	//log.Printf("msgs %v \n", len(chatMessages))
	return &proto.QueryMessagesResponse{
		User1:    request.GetUser1(),
		User2:    request.GetUser2(),
		Messages: chatMessages,
	}, nil
}

func (s Server) UploadFile(server proto.MessagingService_UploadFileServer) error {
	//log.Printf("Is authenticated ? : %v\n", server.Context().Value("isAuthenticated"))
	fileBytes := make([]byte, 0, 1000)
	var fileName string
	var fromUser *proto.User
	var toUser *proto.User
	var chatId string
	var replyChatId string
	var body string
	var timeStamp time.Time
	msg, err := server.Recv()
	if err != nil {
		return status.Error(codes.Unknown, "File info missing")
	}
	fileName = msg.GetInfo().GetFileName()
	fromUser = msg.GetInfo().GetChatMessage().GetFromUser()
	toUser = msg.GetInfo().GetChatMessage().GetToUser()
	chatId = msg.GetInfo().GetChatMessage().GetChatId()
	replyChatId = msg.GetInfo().GetChatMessage().GetReplyForChatId()
	body = msg.GetInfo().GetChatMessage().GetBody()

	if err2 := timeStamp.UnmarshalText([]byte(msg.GetInfo().GetChatMessage().GetTimeStamp())); err2 != nil {
		log.Printf("Trouble while parsing Time; %v\n", err2.Error())
		return status.Error(codes.InvalidArgument, "invalid value for msg time property")
	}

	for {
		msg, err = server.Recv()
		if err != nil {
			if err == io.EOF {
				err2 := server.SendAndClose(&proto.FileUploadStatus{
					Status:  proto.FileUploadStatus_OK,
					Message: fmt.Sprintf("File - %v - uploaded successfully", fileName),
				})
				if err2 != nil {
					return err2
				}
				break
			}
			log.Printf("trouble with client stream %v\n", err.Error())
			return err
		}
		//log.Println("Len " + string(len(msg.GetContent())))
		fileBytes = append(fileBytes, msg.GetContent()...)
	}

	//log.Printf("FileName: %v\n", fileName)
	//log.Printf("Filebytes: %v\n", len(fileBytes))
	docIDs := make([]string, 0, 2)
	chunkSize := 512000
	i := 0
	for {
		lowerLimit := i * chunkSize
		upperLimit := int(math.Min(float64((i+1)*chunkSize), float64(len(fileBytes))))
		chunk := fileBytes[lowerLimit:upperLimit]
		id := gocql.TimeUUID().String()
		if err = session.Query(`INSERT INTO chat_app.documents (id, chat_id, document_name, document_content) VALUES (?, ?, ?, ?)`, id, chatId, fileName, chunk).WithContext(server.Context()).Exec(); err != nil {
			log.Printf("Trouble while inserting doc into db; %v\n", err.Error())
			//return nil, err
		}
		docIDs = append(docIDs, id)
		i++
		if upperLimit == len(fileBytes) {
			break
		}
		//log.Println(len(chunk))
	}
	b, err2 := timeStamp.MarshalText()
	if err2 != nil {
		log.Printf("error while marshaling time; %v\n", err2.Error())
		return status.Error(codes.Internal, "trouble while sending user m")
	}
	_, err = s.BroadcastMessage(server.Context(), &proto.ChatMessage{
		ChatId:         chatId,
		FromUser:       fromUser,
		ToUser:         toUser,
		Body:           body,
		Attachment:     true,
		TimeStamp:      string(b),
		ReplyForChatId: replyChatId,
	})
	if err != nil {
		return err
	}
	if err = session.Query(`UPDATE chat_app.chat_messages SET documents = ? WHERE from_user = ? AND to_user = ? AND time = ?`, docIDs, fromUser.GetUserId(), toUser.GetUserId(), timeStamp).WithContext(server.Context()).Exec(); err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s Server) DownloadFile(message *proto.ChatMessage, server proto.MessagingService_DownloadFileServer) error {

	var (
		documentName string
		fileData     []byte
	)
	scan2 := session.Query(`SELECT document_content, document_name FROM chat_app.documents WHERE chat_id = ? ALLOW FILTERING;`, message.GetChatId()).Iter().Scanner()
	for scan2.Next() {
		var cont []byte
		err2 := scan2.Scan(&cont, &documentName)
		if err2 != nil {
			if err2 == gocql.ErrNotFound {
				return status.Error(codes.InvalidArgument, "invalid query parameters")
			}
			log.Printf("error while querying; %v", err2.Error())
			return status.Error(codes.Internal, "trouble while downloading file")
		}
		//fmt.Printf("Doc cont : %v\n", len(cont))
		fileData = append(fileData, cont...)
	}
	//fmt.Printf("Total Doc cont : %v\n", len(fileData))
	err := server.Send(
		&proto.FileDownloadChunk{
			DownloadData: &proto.FileDownloadChunk_Info{
				Info: &proto.FileDownloadChunk_FileDownload{
					FileName:    fmt.Sprintf("%v - %v", len(fileData), documentName),
					ChatMessage: message,
				},
			},
		})
	if err != nil {
		log.Printf("trouble while downloading file %v\n", err.Error())
		return err
	}
	chunkSize := 1024
	i := 0
	for {
		lowerLimit := i * chunkSize
		upperLimit := int(math.Min(float64((i+1)*chunkSize), float64(len(fileData))))
		chunk := fileData[lowerLimit:upperLimit]
		//fmt.Printf("Sending data %v\n", len(chunk))
		err2 := server.Send(
			&proto.FileDownloadChunk{
				DownloadData: &proto.FileDownloadChunk_Content{
					Content: chunk,
				},
			})
		if err2 != nil {
			fmt.Printf("trouble streaming %v\n", err2.Error())
			return err
		}
		i++
		if upperLimit == len(fileData) {
			break
		}
	}
	return nil
}

var session *gocql.Session

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	session, _ = db.ConnectToDB()
	log.Println("Successfully connected to DB")
	defer session.Close()

	var onlineUsers = make(map[string]*Connection)

	server := &Server{onlineUsers, proto.UnimplementedMessagingServiceServer{}}

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(interceptors.AuthHandlingInterceptor)),
		grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(interceptors.AuthHandlingInterceptor)))
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("error creating the server %v", err)
	}

	//grpcLog.Info("Starting server at port :8080")

	proto.RegisterMessagingServiceServer(grpcServer, server)
	err3 := grpcServer.Serve(listener)
	if err3 != nil {
		log.Fatal("can't serve on specified port")
	}
}
