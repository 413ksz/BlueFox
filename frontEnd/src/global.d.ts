/// <reference types="@solidjs/start/env" />

type Server = {
  id: string;
  name: string;
  avatar: string;
  channels: Channel[];
};

type Channel = {
  id: string;
  name: string;
  messages: Message[];
};

type Message = {
  text: string;
  sender: string;
};

type ActiveChat = {
  type: "friend" | "server";
  mainId: string;
  channelId?: string;
};

type UserProfile = {
  id: string;
  name: string;
  avatar: string;
};

interface User {
  id: string;
  username: string;
  email: string;
  password_hash: string;
  profile_picture_asset_id?: string;
  created_at: Date;
  updated_at?: Date;
  last_online?: Date;
  bio?: string;
  date_of_birth: Date;
  location?: string;
  is_verified: boolean;
}
