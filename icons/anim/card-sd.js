import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const card = el.querySelector('path[d^="M6 22"]');
  animate(card, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const contacts = el.querySelectorAll('path[d^="M8 10"], path[d^="M12 6"], path[d^="M16 6"]');
  animate(contacts, { scaleY: [0, 1.25, 1] }, { duration: τmid, delay: stagger(0.06, { startDelay: τmin }) });
}
