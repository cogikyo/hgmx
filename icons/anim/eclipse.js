import { animate, τmax } from "../motion.js";

export default function (el) {
  const shadow = el.querySelector("path");
  animate(shadow, { rotate: [0, -30, 10, 0] }, { duration: τmax });
}
