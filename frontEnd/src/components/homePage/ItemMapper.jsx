import { For } from "solid-js";

const ItemMapper = ({ items, mounted }) => {
  return (
    <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-6 justify-items-center mb-16">
      <For each={items}>
        {(item) => (
          <a
            href={item.url}
            target="_blank"
            rel="noopener noreferrer"
            className="
              flex flex-col items-center justify-center text-center
              bg-gray-800/70 backdrop-blur-sm rounded-xl
              px-4 py-6 shadow-lg border border-gray-700
              text-gray-300 hover:text-white
              hover:bg-blue-600/30 hover:border-blue-500
              transition-all duration-300 ease-in-out
              transform hover:-translate-y-1 hover:scale-105
              w-full h-32 sm:h-36
            "
            aria-label={`Learn more about ${item.name}`}
          >
            {mounted() && <item.icon className="w-8 h-8 mb-3 text-blue-400" />}
            <span className="text-lg font-semibold">{item.name}</span>
          </a>
        )}
      </For>
    </div>
  );
};

export default ItemMapper;
