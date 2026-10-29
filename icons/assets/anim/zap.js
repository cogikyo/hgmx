import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const bolt = el.querySelector("path");
  animate(bolt, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });
  animate(bolt, { translateX: [3, -1, 0], translateY: [-4, 1, 0] }, { duration: τmax });
}
