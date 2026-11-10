import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const middle = el.querySelectorAll("path");
  animate(middle, { scale: [1, 0.6, 1] }, { duration: τmax });

  const top = el.querySelector('rect[x="2"]');
  animate(top, { translateX: [0, 1.5, 0], translateY: [0, 1.5, 0] }, { duration: τmid });

  const bottom = el.querySelector('rect[x="14"]');
  animate(bottom, { translateX: [0, -1.5, 0], translateY: [0, -1.5, 0] }, { duration: τmid });
}
