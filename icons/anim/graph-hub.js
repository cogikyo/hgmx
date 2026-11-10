import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const hub = el.querySelector('circle[r="4"]');
  animate(hub, { scale: [1, 0.85, 1] }, { duration: τmid });

  const leaves = el.querySelectorAll('circle[r="2"]');
  animate(leaves, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
