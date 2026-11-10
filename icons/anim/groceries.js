import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const bag = el.querySelector('path[d^="M4 12.006"]');
  animate(bag, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const produce = el.querySelectorAll('path[d^="M7 11"], path[d^="M19.424"]');
  animate(produce, { translateY: [-4, 1, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
