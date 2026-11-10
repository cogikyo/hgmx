import { animate, τmax } from "../motion.js";

export default function (el) {
  const rook = el.querySelector("g");
  animate(rook, { translateX: [0, 3, -1, 0] }, { duration: τmax });
}
