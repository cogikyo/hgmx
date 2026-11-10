import { animate, τmax } from "../motion.js";

export default function (el) {
  const whistle = el.querySelector("g");
  animate(whistle, { rotate: [0, -8, 6, -6, 4, 0] }, { duration: τmax });
}
