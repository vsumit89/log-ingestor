import React from "react";
import { sora } from "./navbar.style";

export const Navbar = ({}) => {
  return (
    <nav className="w-full p-4 flex justify-between text-lg shadow-sm items-center">
      <span className={`${sora.className} cursor-pointer`}>
        LogSwift
        <span className="text-xs ml-2">Navigate through your logs swiftly</span>
      </span>
    </nav>
  );
};
