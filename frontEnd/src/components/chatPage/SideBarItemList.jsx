import { For } from "solid-js";
import SideBarItem from "./SideBarItem";
const SideBarItemList = ({ items, isSelectedChat, handleChatSelection }) => {
  return (
    <For each={items()}>
      {(item) => (
        <SideBarItem
          isSelectedChat={isSelectedChat}
          handleChatSelection={handleChatSelection}
          item={item}
        />
      )}
    </For>
  );
};

export default SideBarItemList;
