import { animate, τmax } from "../motion.js";

export default function (el) {
  const upper = el.querySelector('circle[cx="7"]');
  animate(upper, { translateX: [0, 4, -1, 0] }, { duration: τmax });

  const lower = el.querySelector('circle[cx="17"]');
  animate(lower, { translateX: [0, -4, 1, 0] }, { duration: τmax });
}
