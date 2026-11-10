import { animate, τmax, τmin } from "../motion.js";

export default function (el) {
  const moon = el.querySelector('path[d^="M20.985"]');
  animate(moon, { rotate: [0, -30, 10, 0], scale: [1, 0.85, 1] }, { duration: τmax });

  const star = el.querySelectorAll('path[d="M18 5h4"], path[d="M20 3v4"]');
  animate(star, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: τmin });
}
