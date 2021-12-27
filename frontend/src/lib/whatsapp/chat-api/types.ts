export interface Dialog {
    id: string
    name: string
    image: string
    metadata: DialogMetadata | DialogBusinessMetadata | DialogGroupMetadata
    last_time: number // unix timestamp
}

export interface DialogMetadata {
    admins: string[]
    isGroup: false
    participants: string[]
    groupInviteLink: null
    participantsInfo: string[]
}

export interface DialogBusinessMetadata {
    admins: string[]
    isGroup: false
    participants: string[]
    verifiedName: string // "Company Name"",
    businessProfile: {
        email: string // "owner@example.com"
        website: string[] // [ "https://www.example.com" ]
        description: string // "Company description"
    },
    groupInviteLink: null
    participantsInfo: string[]
}

export interface DialogGroupMetadata {
    admins: string[] // [ "5512345678@c.us" ]
    isGroup: true
    announce: string // "all"
    restrict: string // "all"
    participants: string[] // [ "5512345678@c.us" ]
    groupInviteLink: string | null
    participantsInfo: ParticipantInfo[]
}

export interface ParticipantInfo {
    id: string // "5512345678@c.us"
    name: string // "John Doe"
}

export interface Message {
    __rowID: number // not from Chat-API, but our own local database
    ack: string | null // null | "viewed" | "delivered"
    author: string // "5512345678@c.us" | "undefined"
    body: string // "Oi! Teste!"
    caption: string // Used for type="image"
    chatId: string // "5512345678-1234567801@g.us"
    chatName: string // "Chat Name"
    fromMe: boolean // false
    id: string // "false_5512345678-1234567801@g.us_798103AB1C123456780abcdef12345678_5512345678@c.us"
    isForwarded: number // 0
    messageNumber: number // 33
    metadata: MessageMetadataLink | null
    quotedMsgBody: string | null
    quotedMsgId: string | null // "false_5512345678-1234567801@g.us_8AE8545512345678abcdef12345678_5512345678@c.us"
    quotedMsgType: string | null // "image"
    self: number // 0
    senderName: string // "John Doe"
    time: number // 1639597451 (unix timestamp)
    type: string // "chat"
}

export interface MessageMetadataLink {
    linkUrl: string // "https://example.com/abc/",
    linkTitle: string // "About the link",
    linkDescription: string // " "
}

export interface Label {
    id: string
    name: string
    hexColor: string
}
