import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const udder = el.querySelector('path[d^="M2 2"]');
  animate(udder, { scaleY: [1, 0.85, 1.05, 1] }, { duration: τmid });

  const drops = el.querySelectorAll('path[d$="2.8 0"]');
  animate(drops, { translateY: [-4, 1, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
