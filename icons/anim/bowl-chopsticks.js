import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const bowl = el.querySelectorAll('ellipse, path[d^="M22 12"]');
  animate(bowl, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const chopsticks = el.querySelector("g");
  animate(chopsticks, { translateY: [0, -3, 1, 0] }, { duration: τmid, delay: τmin });
}
