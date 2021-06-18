//Generated by the protocol buffer compiler. DO NOT EDIT!
// source: chatapp/chat_app.proto

package proto.kotlin;

@kotlin.jvm.JvmSynthetic
inline fun chatMessage(block: proto.kotlin.ChatMessageKt.Dsl.() -> Unit): proto.kotlin.ChatMessage =
  proto.kotlin.ChatMessageKt.Dsl._create(proto.kotlin.ChatMessage.newBuilder()).apply { block() }._build()
object ChatMessageKt {
  @kotlin.OptIn(com.google.protobuf.kotlin.OnlyForUseByGeneratedProtoCode::class)
  @com.google.protobuf.kotlin.ProtoDslMarker
  class Dsl private constructor(
    @kotlin.jvm.JvmField private val _builder: proto.kotlin.ChatMessage.Builder
  ) {
    companion object {
      @kotlin.jvm.JvmSynthetic
      @kotlin.PublishedApi
      internal fun _create(builder: proto.kotlin.ChatMessage.Builder): Dsl = Dsl(builder)
    }

    @kotlin.jvm.JvmSynthetic
    @kotlin.PublishedApi
    internal fun _build(): proto.kotlin.ChatMessage = _builder.build()

    /**
     * <code>string chat_id = 1;</code>
     */
    var chatId: kotlin.String
      @JvmName("getChatId")
      get() = _builder.getChatId()
      @JvmName("setChatId")
      set(value) {
        _builder.setChatId(value)
      }
    /**
     * <code>string chat_id = 1;</code>
     */
    fun clearChatId() {
      _builder.clearChatId()
    }

    /**
     * <code>.proto.chat_app.User from_user = 2;</code>
     */
    var fromUser: proto.kotlin.User
      @JvmName("getFromUser")
      get() = _builder.getFromUser()
      @JvmName("setFromUser")
      set(value) {
        _builder.setFromUser(value)
      }
    /**
     * <code>.proto.chat_app.User from_user = 2;</code>
     */
    fun clearFromUser() {
      _builder.clearFromUser()
    }
    /**
     * <code>.proto.chat_app.User from_user = 2;</code>
     * @return Whether the fromUser field is set.
     */
    fun hasFromUser(): kotlin.Boolean {
      return _builder.hasFromUser()
    }

    /**
     * <code>.proto.chat_app.User to_user = 3;</code>
     */
    var toUser: proto.kotlin.User
      @JvmName("getToUser")
      get() = _builder.getToUser()
      @JvmName("setToUser")
      set(value) {
        _builder.setToUser(value)
      }
    /**
     * <code>.proto.chat_app.User to_user = 3;</code>
     */
    fun clearToUser() {
      _builder.clearToUser()
    }
    /**
     * <code>.proto.chat_app.User to_user = 3;</code>
     * @return Whether the toUser field is set.
     */
    fun hasToUser(): kotlin.Boolean {
      return _builder.hasToUser()
    }

    /**
     * <code>string body = 4;</code>
     */
    var body: kotlin.String
      @JvmName("getBody")
      get() = _builder.getBody()
      @JvmName("setBody")
      set(value) {
        _builder.setBody(value)
      }
    /**
     * <code>string body = 4;</code>
     */
    fun clearBody() {
      _builder.clearBody()
    }

    /**
     * <code>string time_stamp = 5;</code>
     */
    var timeStamp: kotlin.String
      @JvmName("getTimeStamp")
      get() = _builder.getTimeStamp()
      @JvmName("setTimeStamp")
      set(value) {
        _builder.setTimeStamp(value)
      }
    /**
     * <code>string time_stamp = 5;</code>
     */
    fun clearTimeStamp() {
      _builder.clearTimeStamp()
    }

    /**
     * <code>string reply_for_chat_id = 6;</code>
     */
    var replyForChatId: kotlin.String
      @JvmName("getReplyForChatId")
      get() = _builder.getReplyForChatId()
      @JvmName("setReplyForChatId")
      set(value) {
        _builder.setReplyForChatId(value)
      }
    /**
     * <code>string reply_for_chat_id = 6;</code>
     */
    fun clearReplyForChatId() {
      _builder.clearReplyForChatId()
    }

    /**
     * <code>bool attachment = 7;</code>
     */
    var attachment: kotlin.Boolean
      @JvmName("getAttachment")
      get() = _builder.getAttachment()
      @JvmName("setAttachment")
      set(value) {
        _builder.setAttachment(value)
      }
    /**
     * <code>bool attachment = 7;</code>
     */
    fun clearAttachment() {
      _builder.clearAttachment()
    }
  }
}
@kotlin.jvm.JvmSynthetic
inline fun proto.kotlin.ChatMessage.copy(block: proto.kotlin.ChatMessageKt.Dsl.() -> Unit): proto.kotlin.ChatMessage =
  proto.kotlin.ChatMessageKt.Dsl._create(this.toBuilder()).apply { block() }._build()
