//Generated by the protocol buffer compiler. DO NOT EDIT!
// source: chatapp/chat_app.proto

package proto.kotlin;

@kotlin.jvm.JvmSynthetic
inline fun fileUploadStatus(block: proto.kotlin.FileUploadStatusKt.Dsl.() -> Unit): proto.kotlin.FileUploadStatus =
  proto.kotlin.FileUploadStatusKt.Dsl._create(proto.kotlin.FileUploadStatus.newBuilder()).apply { block() }._build()
object FileUploadStatusKt {
  @kotlin.OptIn(com.google.protobuf.kotlin.OnlyForUseByGeneratedProtoCode::class)
  @com.google.protobuf.kotlin.ProtoDslMarker
  class Dsl private constructor(
    @kotlin.jvm.JvmField private val _builder: proto.kotlin.FileUploadStatus.Builder
  ) {
    companion object {
      @kotlin.jvm.JvmSynthetic
      @kotlin.PublishedApi
      internal fun _create(builder: proto.kotlin.FileUploadStatus.Builder): Dsl = Dsl(builder)
    }

    @kotlin.jvm.JvmSynthetic
    @kotlin.PublishedApi
    internal fun _build(): proto.kotlin.FileUploadStatus = _builder.build()

    /**
     * <code>string message = 1;</code>
     */
    var message: kotlin.String
      @JvmName("getMessage")
      get() = _builder.getMessage()
      @JvmName("setMessage")
      set(value) {
        _builder.setMessage(value)
      }
    /**
     * <code>string message = 1;</code>
     */
    fun clearMessage() {
      _builder.clearMessage()
    }

    /**
     * <code>.proto.chat_app.FileUploadStatus.UploadStatusCode status = 2;</code>
     */
    var status: proto.kotlin.FileUploadStatus.UploadStatusCode
      @JvmName("getStatus")
      get() = _builder.getStatus()
      @JvmName("setStatus")
      set(value) {
        _builder.setStatus(value)
      }
    /**
     * <code>.proto.chat_app.FileUploadStatus.UploadStatusCode status = 2;</code>
     */
    fun clearStatus() {
      _builder.clearStatus()
    }
  }
}
@kotlin.jvm.JvmSynthetic
inline fun proto.kotlin.FileUploadStatus.copy(block: proto.kotlin.FileUploadStatusKt.Dsl.() -> Unit): proto.kotlin.FileUploadStatus =
  proto.kotlin.FileUploadStatusKt.Dsl._create(this.toBuilder()).apply { block() }._build()
