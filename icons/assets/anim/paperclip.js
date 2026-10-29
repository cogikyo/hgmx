import { animate, τmax } from "../motion.js";

export default function (el) {
  const clip = el.querySelector("path");
  animate(clip, { translateX: [0, 3, -1, 0], translateY: [0, -3, 1, 0] }, { duration: τmax });
}
