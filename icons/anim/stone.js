import { animate, τmax } from "../motion.js";

export default function (el) {
  const stone = el.querySelectorAll("path");
  animate(stone, { translateY: [0, -2, 1, 0] }, { duration: τmax });
}
