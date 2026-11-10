import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const walls = el.querySelector('path[d^="M22 19"]');
  animate(walls, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const slats = el.querySelectorAll('path[d="M6 13h12"], path[d="M6 17h12"]');
  animate(slats, { scaleX: [0, 1] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
