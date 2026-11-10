import { animate, τmax } from "../motion.js";

export default function (el) {
  const card = el.querySelector("rect");
  animate(card, { translateX: [-2, 0.5, 0] }, { duration: τmax });
}
