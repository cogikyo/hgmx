import { animate, τmax } from "../motion.js";

export default function (el) {
  const nut = el.querySelector("path");
  animate(nut, { rotate: [0, 360], scale: [1, 0.85, 1] }, { duration: τmax });
}
