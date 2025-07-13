import { Switch, Match } from "solid-js";
import { Motion } from "solid-motionone";

const SideBarItem = ({ isSelectedChat, handleChatSelection, item }) => {
  const itemType = item.username ? "friends" : "servers";

  return (
    <Motion.div
      initial={{ opacity: 0, x: -20 }}
      animate={{ opacity: 1, x: 0 }}
      exit={{ opacity: 0, x: -20 }}
      hover={{ scale: 1.03 }}
      press={{ scale: 0.98 }}
      transition={{ duration: 0.15, ease: "easeOut" }}
    >
      <button
        classList={{
          "w-full flex items-center justify-start gap-4 py-3 px-6 transition-all duration-300 text-left rounded-lg": true,
          "bg-blue-600 text-white shadow-lg": isSelectedChat(item),
          "bg-transparent text-gray-300 hover:bg-gray-700 hover:text-white":
            !isSelectedChat(item),
        }}
        onClick={() => {
          handleChatSelection(itemType, item.id);
        }}
        aria-current={isSelectedChat(item) ? "page" : undefined}
      >
        <Switch>
          <Match when={item.username}>
            <div class="h-10 w-10 rounded-full overflow-hidden flex-shrink-0 border-2 border-gray-600">
              <img
                src={item.profile_picture_asset_id}
                alt={item.username}
                class="h-full w-full object-cover"
              />
            </div>
            <span class="truncate text-lg font-medium">{item.username}</span>{" "}
          </Match>
          <Match when={item.name}>
            <div class="h-10 w-10 rounded-full overflow-hidden flex-shrink-0 border-2 border-gray-600">
              <img
                src={item.avatar}
                alt={item.name}
                class="h-full w-full object-cover"
              />
            </div>
            <span class="truncate text-lg font-medium">{item.name}</span>
          </Match>
        </Switch>
      </button>
    </Motion.div>
  );
};

export default SideBarItem;
