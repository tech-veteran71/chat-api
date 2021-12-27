export type { Message, MessageMetadataLink, Dialog, Label } from './types'

export { getMessages } from './get-messages'
export type { GetMessages, GetMessagesResponse } from './get-messages'

export { getDialogs } from './get-dialogs'
export type { GetDialogs, GetDialogsResponse } from './get-dialogs'

export { getLabels } from './get-labels'
export type { GetLabelsResponse } from './get-labels'

export { sendMessage } from './send-message'
export type { SendMessage, SendMessageResponse } from './send-message'

export { checkPhone } from './check-phone'
export type { CheckPhone, CheckPhoneResponse } from './check-phone'

export { labelChat } from './label-chat'
export type { LabelChat, LabelChatResponse } from './label-chat'

export { unlabelChat } from './unlabel-chat'
export type { UnlabelChat, UnlabelChatResponse } from './unlabel-chat'

export { archiveChat } from './archive-chat'
export type { ArchiveChat, ArchiveChatResponse } from './archive-chat'

export { readChat } from './read-chat'
export type { ReadChat, ReadChatResponse } from './read-chat'

export { unreadChat } from './unread-chat'
export type { UnreadChat, UnreadChatResponse } from './unread-chat'

export { deleteMessage } from './delete-message'
export type { DeleteMessage, DeleteMessageResponse } from './delete-message'
