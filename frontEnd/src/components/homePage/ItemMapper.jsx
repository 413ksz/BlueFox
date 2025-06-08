import { For } from "solid-js"; // Import the For component

const ItemMapper = ({ items }) => {
  return (
    <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-6 justify-items-center">
      <For each={items}>
        {(item) => (
          <a
            href={item.url}
            target="_blank"
            rel="noopener noreferrer"
            className="bg-gray-800/80 backdrop-blur-md rounded-xl px-4 py-2 shadow-md border border-gray-700 text-gray-300 hover:text-white hover:bg-blue-500/20 transition-colors"
            aria-label={`Learn more about ${item.name}`}
          >
            {item.name}
          </a>
        )}
      </For>
    </div>
  );
};

export default ItemMapper;
