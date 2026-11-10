import { animate, stagger, τmid } from "../motion.js";

export default function (el) {
  const links = el.querySelectorAll("path");
  animate(links, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const stops = [el.querySelector('circle[cy="4"]'), el.querySelector('circle[cx="4"]'), el.querySelector('circle[cx="20"]'), el.querySelector('circle[cy="20"]')];
  animate(stops, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.09) });
}
