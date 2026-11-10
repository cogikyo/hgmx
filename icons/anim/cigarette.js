import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const ember = el.querySelector('path[d^="M21 16"]');
  animate(ember, { scale: [1, 1.25, 1] }, { duration: τmid });

  const smoke = el.querySelectorAll('path[d^="M18 8"], path[d^="M22 8"]');
  animate(smoke, { translateY: [3, -1, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
