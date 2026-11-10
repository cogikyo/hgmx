import { animate, τmax } from "../motion.js";

export default function (el) {
  const top = el.querySelector('rect[x="3"]');
  animate(top, { translateX: [2, -0.5, 0], translateY: [2, -0.5, 0] }, { duration: τmax });

  const bottom = el.querySelector('rect[x="11"]');
  animate(bottom, { translateX: [-2, 0.5, 0], translateY: [-2, 0.5, 0] }, { duration: τmax });
}
