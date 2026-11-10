import { animate, τmax } from "../motion.js";

export default function (el) {
  const moon = el.querySelector("path");
  animate(moon, { rotate: [0, -30, 10, 0], scale: [1, 0.85, 1] }, { duration: τmax });
}
