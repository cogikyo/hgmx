import { animate, τmax } from "../motion.js";

export default function (el) {
  const reels = el.querySelectorAll("circle");
  animate(reels, { rotate: [0, 360], pathLength: [0.25, 1] }, { duration: τmax });
}
