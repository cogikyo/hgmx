import { animate, τmax } from "../motion.js";

export default function (el) {
  const cloth = el.querySelector("path");
  animate(cloth, { rotate: [0, -10, 6, 0], scale: [1, 0.9, 1] }, { duration: τmax });
}
