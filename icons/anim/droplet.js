import { animate, τmax } from "../motion.js";

export default function (el) {
  const drop = el.querySelector("path");
  animate(drop, { translateY: [-5, 1, 0], scaleY: [1.15, 0.85, 1], scaleX: [0.85, 1.15, 1] }, { duration: τmax });
}
