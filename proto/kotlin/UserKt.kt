//Generated by the protocol buffer compiler. DO NOT EDIT!
// source: chatapp/chat_app.proto

package proto.kotlin;

@kotlin.jvm.JvmSynthetic
inline fun user(block: proto.kotlin.UserKt.Dsl.() -> Unit): proto.kotlin.User =
  proto.kotlin.UserKt.Dsl._create(proto.kotlin.User.newBuilder()).apply { block() }._build()
object UserKt {
  @kotlin.OptIn(com.google.protobuf.kotlin.OnlyForUseByGeneratedProtoCode::class)
  @com.google.protobuf.kotlin.ProtoDslMarker
  class Dsl private constructor(
    @kotlin.jvm.JvmField private val _builder: proto.kotlin.User.Builder
  ) {
    companion object {
      @kotlin.jvm.JvmSynthetic
      @kotlin.PublishedApi
      internal fun _create(builder: proto.kotlin.User.Builder): Dsl = Dsl(builder)
    }

    @kotlin.jvm.JvmSynthetic
    @kotlin.PublishedApi
    internal fun _build(): proto.kotlin.User = _builder.build()

    /**
     * <code>string user_id = 1;</code>
     */
    var userId: kotlin.String
      @JvmName("getUserId")
      get() = _builder.getUserId()
      @JvmName("setUserId")
      set(value) {
        _builder.setUserId(value)
      }
    /**
     * <code>string user_id = 1;</code>
     */
    fun clearUserId() {
      _builder.clearUserId()
    }

    /**
     * <code>string name = 2;</code>
     */
    var name: kotlin.String
      @JvmName("getName")
      get() = _builder.getName()
      @JvmName("setName")
      set(value) {
        _builder.setName(value)
      }
    /**
     * <code>string name = 2;</code>
     */
    fun clearName() {
      _builder.clearName()
    }

    /**
     * <code>.proto.chat_app.User.UserStatusCode status = 3;</code>
     */
    var status: proto.kotlin.User.UserStatusCode
      @JvmName("getStatus")
      get() = _builder.getStatus()
      @JvmName("setStatus")
      set(value) {
        _builder.setStatus(value)
      }
    /**
     * <code>.proto.chat_app.User.UserStatusCode status = 3;</code>
     */
    fun clearStatus() {
      _builder.clearStatus()
    }

    /**
     * <code>.proto.chat_app.User.Type type = 4;</code>
     */
    var type: proto.kotlin.User.Type
      @JvmName("getType")
      get() = _builder.getType()
      @JvmName("setType")
      set(value) {
        _builder.setType(value)
      }
    /**
     * <code>.proto.chat_app.User.Type type = 4;</code>
     */
    fun clearType() {
      _builder.clearType()
    }
  }
}
@kotlin.jvm.JvmSynthetic
inline fun proto.kotlin.User.copy(block: proto.kotlin.UserKt.Dsl.() -> Unit): proto.kotlin.User =
  proto.kotlin.UserKt.Dsl._create(this.toBuilder()).apply { block() }._build()
