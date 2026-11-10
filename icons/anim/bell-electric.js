import { animate, τmax, τmin } from "../motion.js";

export default function (el) {
  const gong = el.querySelectorAll('circle[r="7"], path[d^="M9 9"]');
  animate(gong, { translateX: [0, -1, 1, -0.5, 0] }, { duration: τmax });
  const sound = el.querySelector('path[d^="M18.8 4"]');
  animate(sound, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: τmin });
}
