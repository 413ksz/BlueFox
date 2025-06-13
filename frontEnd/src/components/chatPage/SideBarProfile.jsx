const SideBarProfile = ({ currentUser }) => {
  return (
    <div class="p-4 border-t border-gray-700 flex items-center bg-gray-700">
      <img
        src={currentUser()?.profile_picture_asset_id}
        alt="User Avatar"
        class="w-10 h-10 rounded-full mr-3"
      />
      <span class="font-semibold text-white">{currentUser()?.username}</span>
    </div>
  );
};

export default SideBarProfile;
