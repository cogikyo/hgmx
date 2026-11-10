import { animate, τmax } from "../motion.js";

export default function (el) {
  const egg = el.querySelector("path");
  animate(egg, { rotate: [0, -15, 10, -5, 0] }, { duration: τmax });
}
