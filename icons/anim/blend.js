import { animate, τmax } from "../motion.js";

export default function (el) {
  const upper = el.querySelector('circle[cx="15"]');
  animate(upper, { translateX: [3, -0.5, 0], translateY: [-3, 0.5, 0] }, { duration: τmax });

  const lower = el.querySelector('circle[cx="9"]');
  animate(lower, { translateX: [-3, 0.5, 0], translateY: [3, -0.5, 0] }, { duration: τmax });
}
