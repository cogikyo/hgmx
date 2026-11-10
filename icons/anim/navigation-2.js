import { animate, τmax } from "../motion.js";

export default function (el) {
  const arrow = el.querySelector("polygon");
  animate(arrow, { translateY: [5, -2, 0], rotate: [-30, 10, 0], opacity: [0, 1] }, { duration: τmax });
}
