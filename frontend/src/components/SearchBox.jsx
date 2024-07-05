import React from "react";

const SearchBox = ({ onChange, setAuthors }) => {
  const handleChange = (e) => {
    e.preventDefault();
    if (e.target.value !== "") {
      onChange(e.target.value);
    } else {
      setAuthors([]);
    }
  };

  return (
    <form className="w-full">
      <input
        type="text"
        onChange={handleChange}
        placeholder="Search author name"
        className="w-full border border-black rounded-full p-3 border-opacity-20 bg-searchBoxFill bg-opacity-20 text-xs font-bold pl-5"
      />
    </form>
  );
};

export default SearchBox;
