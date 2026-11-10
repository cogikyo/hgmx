import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const port = el.querySelector('path[d^="M22 9"]');
  animate(port, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const contacts = el.querySelector('path[d="M8 12h8"]');
  animate(contacts, { scaleX: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
