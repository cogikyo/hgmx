import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const line = el.querySelector("path");
  animate(line, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const tip = el.querySelector('circle[cx="18"]');
  animate(tip, { translateX: [-8, 1, 0], translateY: [8, -1, 0], opacity: [0, 1] }, { duration: τmax });
}
