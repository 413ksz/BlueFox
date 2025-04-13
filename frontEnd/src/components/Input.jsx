import { Motion } from "solid-motionone";

const Input = ({
  mounted,
  context,
  placeHolder,
  IconName,
  value,
  setValue,
  type,
}) => {
  return (
    <Motion.div
      class="space-y-2"
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
    >
      <label for="name" class="text-gray-300 flex items-center gap-2">
        {mounted() && <IconName class="w-6 h-6" />}
        <span class="font-medium">{context}</span>
      </label>
      <input
        id={context}
        type={type}
        value={value()}
        onChange={(e) => setValue(e.target.value)}
        placeholder={placeHolder}
        class="bg-gray-800 text-white border-gray-700 placeholder:text-gray-400 rounded-md px-4 py-2 w-full focus:outline-none focus:ring-2 focus:ring-blue-500"
        autocomplete="name"
      />
    </Motion.div>
  );
};

export default Input;
