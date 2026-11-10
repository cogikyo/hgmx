import { animate, τmax } from "../motion.js";

export default function (el) {
  const card = el.querySelector("rect");
  animate(card, { translateX: [0, -2, 1, 0] }, { duration: τmax });
}
