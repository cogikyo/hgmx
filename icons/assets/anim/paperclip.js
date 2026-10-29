import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const clip = el.querySelector("path");
  animate(clip, { translateX: [4, -1, 0], translateY: [-4, 1, 0] }, { duration: τmax });
  animate(clip, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });
}
