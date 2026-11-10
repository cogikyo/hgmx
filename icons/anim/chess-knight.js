import { animate, τmax } from "../motion.js";

export default function (el) {
  const knight = el.querySelector("g");
  animate(knight, { translateY: [0, -3, 1, 0], rotate: [0, -12, 4, 0] }, { duration: τmax });
}
