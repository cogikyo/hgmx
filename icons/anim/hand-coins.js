import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const hand = el.querySelectorAll("path");
  animate(hand, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const coins = [el.querySelector('circle[cx="6"]'), el.querySelector('circle[cx="16"]')];
  animate(coins, { translateY: [-4, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: stagger(0.12, { startDelay: τmin }) });
}
