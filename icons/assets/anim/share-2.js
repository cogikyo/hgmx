import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const source = el.querySelector('circle[cx="6"]');
  animate(source, { scale: [1, 0.85, 1] }, { duration: τmid });

  const peers = el.querySelectorAll('circle[cx="18"]');
  animate(peers, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.09, { startDelay: τmin }) });
}
