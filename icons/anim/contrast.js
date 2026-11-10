import { animate, τmax } from "../motion.js";

export default function (el) {
  const half = el.querySelector("path");
  animate(half, { rotate: [0, 360] }, { duration: τmax });
}
