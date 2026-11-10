import { animate, τmax } from "../motion.js";

export default function (el) {
  const card = el.querySelectorAll("rect, path");
  animate(card, { rotateY: [0, 360] }, { duration: τmax });
}
