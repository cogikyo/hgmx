import { animate, τmax } from "../motion.js";

export default function (el) {
  const backpack = el.querySelector("g");
  animate(backpack, { rotate: [0, -10, 6, -3, 0] }, { duration: τmax });
}
