package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/gocql/gocql"
	"github.com/webdevgopi/chatApp-gRPC/proto"
	"google.golang.org/grpc"
)

var client proto.MessagingServiceClient
var wait *sync.WaitGroup

func init() {
	wait = &sync.WaitGroup{}
}

func connect(user *proto.ConnectUserRequest, ctx context.Context, opt grpc.CallOption) error {
	var streamerror error

	stream, err := client.ConnectUser(ctx, user, opt)

	if err != nil {
		return fmt.Errorf("connection failed: %v\n", err)
	}

	wait.Add(1)
	go func(str proto.MessagingService_ConnectUserClient, ctx context.Context, opt grpc.CallOption) {
		defer wait.Done()

		for {
			msg, err2 := str.Recv()
			if err2 != nil {
				if err2 == io.EOF {
					streamerror = fmt.Errorf("User logged out\n")
				}
				streamerror = fmt.Errorf("Error reading message: %v\n", err2)
				return
			} else {
				errChannel := make(chan error)
				go func(msg *proto.ChatMessage, ctx context.Context, opt grpc.CallOption, errChannel chan error) {
					fmt.Printf("updating msg status \n")
					res, err3 := client.UpdateMessageStatus(ctx, &proto.ChatMessageStatus{
						ChatMessage: msg,
						MsgStatus:   proto.ChatMessageStatus_SEEN,
					}, opt)
					errChannel <- err3
					if err3 != nil {
						fmt.Printf("trouble while updating msg status %v\n", err3.Error())
					}
					fmt.Printf("res : %v\n", res)
				}(msg, ctx, opt, errChannel)
				<-errChannel
				close(errChannel)
			}

			fmt.Printf("ChatId : %v,\nFromUser: %v,\nBody: %v,\nTime: %v,\n", msg.GetChatId(), msg.GetFromUser(), msg.GetBody(), msg.GetTimeStamp())
		}
	}(stream, ctx, opt)
	//wait.Wait()
	fmt.Println("exiting connect func")
	//fmt.Printf(streamerror.Error())
	return streamerror
}

func main() {
	//timestamp := time.Now()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	done := make(chan int)

	var id = flag.String("ID", "", "The id of the user")
	flag.Parse()
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldnt connect to service: %v", err)
	}

	client = proto.NewMessagingServiceClient(conn)

	md := metadata.Pairs("authorization", "bearer i_am_an_auth_token")
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	// Make RPC using the context with the metadata.
	var header metadata.MD

	//header := metadata.New(
	//	map[string]string{
	//		"authorization": "bearer i_am_an_auth_token",
	//	})
	//header.Set("bearer ", "i_am_an_auth_token")
	//ctx := context.WithValue(context.Background(), "authorization", "bearer i_am_an_auth_token")

	user, err := client.GetUser(ctx, &proto.QueryUser{
		UserId: *id,
	})
	if err != nil {
		fmt.Printf("Trouble while fetching user; %v", err.Error())
		return
	}
	fmt.Printf("Userobj %v\n", user.String())
	userType := proto.ConnectUserRequest_UNKNOWN
	if user.GetType() == proto.User_EXPERT {
		userType = proto.ConnectUserRequest_EXPERT
	}
	if user.GetType() == proto.User_PATIENT {
		userType = proto.ConnectUserRequest_PATIENT
	}
	connectUserObj := &proto.ConnectUserRequest{
		UserId: user.GetUserId(),
		Name:   user.GetName(),
		Device: fmt.Sprintf("%v-laptop", user.GetName()),
		Type:   userType,
	}
	errr := connect(connectUserObj, ctx, grpc.Header(&header))
	if errr != nil {
		return
	}
	fmt.Printf("Client-Streaming\n")
	stream, _ := client.UploadFile(ctx, grpc.Header(&header))
	tempChatId := gocql.TimeUUID().String()
	ts, err3 := time.Now().MarshalText()
	if err3 != nil {
		fmt.Printf("error while marshaling time; %v\n", err3.Error())
	}
	err2 := stream.Send(&proto.FileUploadChunk{
		UploadData: &proto.FileUploadChunk_Info{
			Info: &proto.FileUploadChunk_FileUpload{
				FileName: "a.pdf",
				ChatMessage: &proto.ChatMessage{
					ChatId:     tempChatId,
					FromUser:   user,
					ToUser:     user,
					Body:       "file from u to u",
					TimeStamp:  string(ts),
					Attachment: true,
				},
			},
		},
	})
	if err2 != nil {
		fmt.Printf("Trouble while client streaming %v\n", err2.Error())
		return
	}
	const BufferSize = 1024
	file, errt := os.Open("C:\\Users\\gopia\\Desktop\\ERP-IIT BBS-SREEK.pdf")
	if errt != nil {
		fmt.Println(errt)
		return
	}
	defer file.Close()
	r := bufio.NewReader(file)
	for {
		buffer := make([]byte, BufferSize)
		bytesread, erry := r.Read(buffer)
		if erry != nil {
			if err2 == io.EOF {
				errn := stream.CloseSend()
				if errn != nil {
					return
				}
			}
			break
		}
		err2 = stream.Send(&proto.FileUploadChunk{
			UploadData: &proto.FileUploadChunk_Content{
				Content: buffer[0:bytesread],
			},
		})
		if err2 != nil {
			fmt.Printf("Trouble while client streaming %v\n", err2.Error())
			return
		}
		fmt.Println("bytes read: ", bytesread)
	}
	streamRes, erru := stream.CloseAndRecv()
	if erru != nil {
		fmt.Printf("error %v", erru.Error())
	}
	fmt.Printf("Stream Res : %v\n", streamRes)

	fileDownload, err := client.DownloadFile(ctx, &proto.ChatMessage{
		ChatId: tempChatId,
	}, grpc.Header(&header))
	if err != nil {
		fmt.Printf("error1 %v", err.Error())
		return
	}
	recv, errn := fileDownload.Recv()
	if errn != nil {
		if errn == io.EOF {
			fmt.Printf("didn't even get file info")
		}
		fmt.Printf("error2 %v", errn.Error())
		return
	}
	fileName := recv.GetInfo().GetFileName()
	//tempChatMessage := recv.GetInfo().GetChatMessage()

	fileCont := make([]byte, 0, 1024)
	for {
		recv, errn = fileDownload.Recv()
		if errn != nil {
			if errn == io.EOF {
				fmt.Println("File downloaded successfully")
			} else {
				fmt.Printf("error3 %v", errn.Error())
				return
			}
			break
		}
		fileCont = append(fileCont, recv.GetContent()...)
	}
	err = ioutil.WriteFile(fileName, fileCont, 0644)
	if err != nil {
		log.Printf("Trouble while writing to file %v\n", err.Error())
	}

	wait.Add(1)
	go func(user *proto.User, ctx context.Context, opt grpc.CallOption) {
		defer wait.Done()
		fmt.Println("Enter the user id (with whom you're willing to chat) :")
		reader := bufio.NewReader(os.Stdin)
		chatUserId, _ := reader.ReadString('\n')
		chatUserId = strings.TrimSpace(chatUserId)
		input := &proto.QueryUser{
			UserId: chatUserId,
		}
		res, err2 := client.GetUser(ctx, input, opt)
		if err2 != nil {
			fmt.Printf("Trouble while fetching user: %v => %v\n", chatUserId, err2.Error())
			return
		}
		fmt.Printf("User status %v\n", res.GetStatus().String())
		fmt.Printf("You can send messages to user %v\n", res.GetName())
		b, err2 := time.Now().MarshalText()
		if err2 != nil {
			fmt.Printf("error while marshaling time; %v\n", err2.Error())
		}
		fmt.Printf("Getting all messages\n")
		msgs, err2 := client.GetMessages(ctx, &proto.QueryMessagesRequest{
			User1:          user,
			User2:          res,
			TimeConstraint: string(b),
			Limit:          20,
		}, opt)
		if err2 != nil {
			fmt.Printf("trouble while getting msgs %v\n", err2.Error())
			return
		}
		for i, m := range msgs.GetMessages() {
			fmt.Printf("msg %v : %v\n", i, m)
		}
		fmt.Println("Type message")

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			b, err3 := time.Now().MarshalText()
			if err3 != nil {
				fmt.Printf("error while marshaling time; %v\n", err3.Error())
			}
			msg := &proto.ChatMessage{
				ChatId:    gocql.TimeUUID().String(),
				FromUser:  user,
				ToUser:    res,
				Body:      scanner.Text(),
				TimeStamp: string(b),
			}

			msgRes, err4 := client.BroadcastMessage(ctx, msg, opt)
			if err4 != nil {
				fmt.Printf("Error Sending Message: %v\n", err4)
				break
			}
			errChannel := make(chan error)
			go func(msg *proto.ChatMessage, ctx context.Context, opt grpc.CallOption, errChannel chan error) {
				msgres, err5 := client.GetMessageStatus(ctx, msg)
				errChannel <- err5
				if err5 != nil {
					fmt.Printf("trouble while getting msg status %v\n", err5.Error())
				}
				fmt.Printf("res : %v\n", msgres)
			}(msg, ctx, grpc.Header(&header), errChannel)
			<-errChannel
			close(errChannel)
			fmt.Printf("Message Response for chatId: %v => %v\n", msgRes.GetChatMessage().GetChatId(), msgRes.GetMsgStatus())
		}

	}(user, ctx, grpc.Header(&header))

	go func() {
		wait.Wait()
		close(done)
	}()

	<-done
}
