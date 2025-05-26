import { TbChevronDown, TbChevronUp, TbServer, TbUser } from "solid-icons/tb";
import { createSignal, createEffect, onMount, For, Show } from "solid-js";
import { Motion, Presence } from "solid-motionone";
import SideBarToggle from "~/components/chatPage/SideBarButton";
import SideBarItem from "~/components/chatPage/SideBarItem";
import SideBarItemList from "~/components/chatPage/SideBarItemList";
import SideBarProfile from "~/components/chatPage/SideBarProfile";

const ChatPage = () => {
  const [servers, setServers] = createSignal<Server[]>([
    {
      id: "s1",
      name: "Gaming Zone",
      avatar: "https://source.unsplash.com/random/100x100/?game&1",
      channels: [
        {
          id: "c1",
          name: "General",
          messages: [{ text: "Welcome to General!", sender: "other" }],
        },
        { id: "c2", name: "Off-Topic", messages: [] },
        { id: "c3", name: "LFG", messages: [] },
      ],
    },
    {
      id: "s2",
      name: "Study Group",
      avatar: "https://source.unsplash.com/random/100x100/?book&1",
      channels: [
        { id: "c4", name: "Math", messages: [] },
        { id: "c5", name: "Science", messages: [] },
        { id: "c6", name: "Literature", messages: [] },
      ],
    },
    {
      id: "s3",
      name: "Music Lovers",
      avatar: "https://source.unsplash.com/random/100x100/?music&1",
      channels: [
        { id: "c7", name: "General", messages: [] },
        { id: "c8", name: "Jazz", messages: [] },
        { id: "c9", name: "Rock", messages: [] },
      ],
    },
  ]);

  const [friends, setFriends] = createSignal<User[]>([
    {
      id: "lkw8r1h2g3j9k4m0p",
      username: "user_7654",
      email: "user7654@example.com",
      password_hash: "a_secure_hash_would_go_here",
      profile_picture_asset_id: "asset_abc123d",
      created_at: new Date("2024-03-10T10:30:00.000Z"),
      updated_at: new Date("2024-05-15T14:45:00.000Z"),
      last_online: new Date("2025-05-23T22:15:00.000Z"),
      bio: "Just a dummy user for testing purposes, interested in tech.",
      date_of_birth: new Date("1995-11-20T00:00:00.000Z"),
      location: "New York, NY",
      is_verified: true,
    },
    {
      id: "xqy5z2v6b7n8m9c0x",
      username: "test_1238",
      email: "test1238@example.com",
      password_hash: "a_secure_hash_would_go_here",
      profile_picture_asset_id: undefined,
      created_at: new Date("2023-01-25T08:00:00.000Z"),
      updated_at: undefined,
      last_online: new Date("2025-05-22T19:00:00.000Z"),
      bio: "Just a dummy user for testing purposes, interested in music.",
      date_of_birth: new Date("1988-04-12T00:00:00.000Z"),
      location: "Paris, France",
      is_verified: false,
    },
    {
      id: "pjm9o0i1u2y3t4r5e",
      username: "dev_9012",
      email: "dev9012@example.com",
      password_hash: "a_secure_hash_would_go_here",
      profile_picture_asset_id: "asset_def456g",
      created_at: new Date("2024-10-01T15:00:00.000Z"),
      updated_at: new Date("2025-01-05T11:00:00.000Z"),
      last_online: undefined,
      bio: undefined,
      date_of_birth: new Date("2000-07-07T00:00:00.000Z"),
      location: undefined,
      is_verified: true,
    },
  ]);

  const [showFriends, setShowFriends] = createSignal(false);
  const [showServers, setShowServers] = createSignal(false);
  const [iconsLoaded, setIconsLoaded] = createSignal(false);
  const [currentUser, setCurrentUser] = createSignal<User | null>({
    id: "pjm9o0i1u2y3t4r5e",
    username: "dev_9012",
    email: "dev9012@example.com",
    password_hash: "a_secure_hash_would_go_here",
    profile_picture_asset_id: "asset_def456g",
    created_at: new Date("2024-10-01T15:00:00.000Z"),
    updated_at: new Date("2025-01-05T11:00:00.000Z"),
    last_online: undefined,
    bio: undefined,
    date_of_birth: new Date("2000-07-07T00:00:00.000Z"),
    location: undefined,
    is_verified: true,
  });

  const [messages, setMessages] = createSignal<Message[]>([]);
  const [newMessage, setNewMessage] = createSignal("");
  let chatContainerRef: HTMLDivElement | undefined;

  const [selectedChat, setSelectedChat] = createSignal<ActiveChat | null>({
    type: "friend",
    mainId: "1",
  });

  const [selectedServerChannels, setSelectedServerChannels] = createSignal<
    Channel[] | undefined
  >(
    selectedChat()?.type === "server"
      ? servers().find((s) => s.id === selectedChat()?.mainId)?.channels
      : undefined
  );
  // Load icons on mount
  onMount(() => {
    setIconsLoaded(true);
  });

  createEffect(() => {
    // Scroll to the bottom whenever messages change
    if (chatContainerRef) {
      chatContainerRef.scrollTop = chatContainerRef.scrollHeight;
    }
  });

  // Load initial messages
  createEffect(() => {
    if (selectedChat()?.type === "server") {
      const server = servers().find((s) => s.id === selectedChat()?.mainId);
      const channel = server?.channels.find(
        (c) => c.id === selectedChat()?.channelId
      );
      setMessages(channel?.messages || []);
    } else if (selectedChat()?.type === "friend") {
      setMessages([
        {
          text: `Chat with ${
            friends().find((f) => f.id === selectedChat()?.mainId)?.username
          }`,
          sender: "other",
        },
        { text: "Hello!", sender: "user" },
      ]);
    }
  });

  const handleSendMessage = () => {
    if (newMessage().trim()) {
      const message = { text: newMessage(), sender: "user" };
      setMessages((prev) => [...prev, message]);
      setNewMessage("");

      // Update the messages in the server data (for demo purposes)
      if (
        selectedChat()?.type === "server" &&
        selectedChat()?.mainId &&
        selectedChat()?.channelId
      ) {
        const serverIndex = servers().findIndex(
          (s) => s.id === selectedChat()?.mainId
        );
        const channelIndex = servers()[serverIndex].channels.findIndex(
          (c) => c.id === selectedChat()?.channelId
        );
        if (channelIndex !== -1) {
          servers()[serverIndex].channels[channelIndex].messages.push(message);
        }
      }
    }
  };

  const handleInputChange = (event: Event) => {
    setNewMessage((event.target as HTMLInputElement).value);
  };

  const handleChatSelection = (
    type: "friends" | "servers",
    id: string,
    channelId?: string
  ) => {
    console.log(selectedServerChannels());
    if (type === "friends") {
      setSelectedChat({ type: "friend", mainId: id });
      // Load messages for the selected friend (replace with your data fetching)
      setMessages([
        {
          text: `Chat with ${friends().find((f) => f.id === id)?.username}`,
          sender: "other",
        },
        { text: "Hello!", sender: "user" },
      ]);
    } else {
      // If it is a server, set the server, and the first channel.
      setSelectedChat({ type: "friend", mainId: id, channelId: channelId });
      if (selectedChat()?.mainId && selectedChat()?.channelId) {
        const server = servers().find((s) => s.id === selectedChat()?.mainId);
        const channel = server?.channels.find(
          (c) => c.id === selectedChat()?.channelId
        );
        setMessages(channel?.messages || []);
      }
    }
  };

  function isSelectedChat(item: Server | UserProfile) {
    const currentChat = selectedChat();
    if (currentChat?.type === "server" && currentChat?.channelId === item.id) {
      return true;
    }

    if (currentChat?.type === "friend" && currentChat?.mainId === item.id) {
      return true;
    }

    return false;
  }

  return (
    <div class="flex h-screen bg-gray-900 text-white">
      {/* Sidebar */}
      <aside class="w-64 bg-gray-800 border-r border-gray-700 flex flex-col">
        <div class="py-4 px-6 border-b border-gray-700">
          <h2 class="text-xl font-semibold">Chats</h2>
        </div>
        <div class="flex-1 overflow-y-auto">
          <div class="py-2">
            <SideBarToggle
              title="Friends"
              IconName={TbUser}
              IconChevronUp={TbChevronUp}
              IconChevronDown={TbChevronDown}
              iconsLoaded={iconsLoaded}
              isOpen={showFriends}
              setIsOpen={setShowFriends}
            />

            <Presence>
              <Show when={showFriends()}>
                <Motion.div
                  initial={{ height: 0, opacity: 0 }}
                  animate={{
                    height: "auto",
                    opacity: 1,
                    transition: {
                      height: {
                        duration: 0.3,
                        easing: [0.17, 0.67, 0.83, 0.67],
                      },
                      opacity: { duration: 0.2 },
                    },
                  }}
                  exit={{
                    height: 0,
                    opacity: 0,
                    transition: {
                      height: { duration: 0.2 },
                      opacity: { duration: 0.2 },
                    },
                  }}
                  class="space-y-1"
                >
                  <SideBarItemList
                    items={friends}
                    isSelectedChat={isSelectedChat}
                    handleChatSelection={handleChatSelection}
                  />
                </Motion.div>
              </Show>
            </Presence>
          </div>
          <div class="py-2">
            <SideBarToggle
              title="Servers"
              IconName={TbServer}
              IconChevronUp={TbChevronUp}
              IconChevronDown={TbChevronDown}
              iconsLoaded={iconsLoaded}
              isOpen={showServers}
              setIsOpen={setShowServers}
            />
            <Presence>
              <Show when={showServers()}>
                <Motion.div
                  initial={{ height: 0, opacity: 0 }}
                  animate={{
                    height: "auto",
                    opacity: 1,
                    transition: {
                      height: {
                        duration: 0.3,
                        easing: [0.17, 0.67, 0.83, 0.67],
                      },
                      opacity: { duration: 0.2 },
                    },
                  }}
                  exit={{
                    height: 0,
                    opacity: 0,
                    transition: {
                      height: { duration: 0.2 },
                      opacity: { duration: 0.2 },
                    },
                  }}
                  class="space-y-1"
                >
                  <SideBarItemList
                    items={servers}
                    isSelectedChat={isSelectedChat}
                    handleChatSelection={handleChatSelection}
                  />
                </Motion.div>
              </Show>
            </Presence>
          </div>
        </div>
        <SideBarProfile currentUser={currentUser} />
      </aside>

      {/* Main Chat Area */}
      <main class="flex-1 flex flex-col">
        {/* Header */}
        <header class="bg-gray-800 py-4 px-6 border-b border-gray-700">
          <div class="container mx-auto">
            <h2 class="text-xl font-semibold">Chat Application</h2>{" "}
            {/* Static Title */}
          </div>
        </header>
        {/* Message Area */}
        <div ref={chatContainerRef} class="flex-grow overflow-y-auto p-4">
          <For each={messages()}>
            {(msg, index) => (
              <div
                classList={{
                  "mb-2 p-3 rounded-lg": true,
                  "bg-blue-600 text-white self-end max-w-[70%]":
                    msg.sender === "user",
                  "bg-gray-700 text-gray-300 self-start max-w-[70%]":
                    msg.sender !== "user",
                }}
              >
                {msg.text}
              </div>
            )}
          </For>
        </div>

        {/* Input Area */}
        <div class="bg-gray-800 py-4 px-6 border-t border-gray-700">
          <div class="container mx-auto flex items-center">
            <input
              type="text"
              class="flex-grow bg-gray-700 text-white rounded-full py-2 px-4 focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Type your message..."
              value={newMessage()}
              onInput={handleInputChange}
              onKeyPress={(event) => {
                if (event.key === "Enter") {
                  handleSendMessage();
                }
              }}
            />
            <button
              class="ml-4 bg-blue-500 hover:bg-blue-600 text-white rounded-full py-2 px-4 font-semibold focus:outline-none focus:ring-2 focus:ring-blue-500"
              onClick={handleSendMessage}
            >
              Send
            </button>
          </div>
        </div>
      </main>
    </div>
  );
};

export default ChatPage;
