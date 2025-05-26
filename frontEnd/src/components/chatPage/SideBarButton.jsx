import { Show } from "solid-js";

const SideBarToggle = ({
  title,
  IconChevronUp,
  IconChevronDown,
  IconName,
  isOpen,
  setIsOpen,
  iconsLoaded,
}) => {
  return (
    <button
      class="w-full flex items-center justify-between py-3 px-6 text-gray-300 hover:text-white hover:bg-gray-700 transition-colors duration-200 rounded-lg group" // Increased padding, more rounded, added group for hover effects
      onClick={() => setIsOpen(!isOpen())}
      aria-expanded={isOpen()}
    >
      <h3 class="font-bold flex items-center gap-3 text-lg">
        <Show when={iconsLoaded()}>
          {
            <IconName class="w-6 h-6 text-blue-400 group-hover:text-blue-300 transition-colors duration-200" />
          }
        </Show>
        {title}
      </h3>
      <Show when={iconsLoaded()}>
        {isOpen() ? (
          <IconChevronUp class="h-6 w-6 text-gray-400 group-hover:text-white transition-colors duration-200" />
        ) : (
          <IconChevronDown class="h-6 w-6 text-gray-400 group-hover:text-white transition-colors duration-200" />
        )}
      </Show>
    </button>
  );
};

export default SideBarToggle;
