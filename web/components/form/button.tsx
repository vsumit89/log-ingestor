import React from "react";

interface ButtonProps {
  text: string;
  onClick: () => void;
}

export const Button: React.FC<ButtonProps> = ({ text, onClick }) => {
  return (
    <button
      className="bg-[#8478E2] text-white px-3 py-1 rounded-sm max-w-[8rem]"
      onClick={onClick}
    >
      {text}
    </button>
  );
};
